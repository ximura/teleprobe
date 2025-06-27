### Database schema 
* rooms: List of all rooms.
* sensors: List of all sensors and what they measure (V or R), linked to rooms.
* measurements: Stores all measurements from sensors (value + timestamp).

```sql
-- Table: rooms
CREATE TABLE rooms (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Table: sensors
CREATE TABLE sensors (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    room_id INTEGER NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    parameter_type TEXT NOT NULL CHECK (parameter_type IN ('V', 'R'))
);

-- Table: measurements
CREATE TABLE measurements (
    id SERIAL PRIMARY KEY,
    sensor_id INTEGER NOT NULL REFERENCES sensors(id) ON DELETE CASCADE,
    timestamp TIMESTAMPTZ NOT NULL,
    value DOUBLE PRECISION NOT NULL
);

-- Index for quick sensor lookup by room
CREATE INDEX idx_sensors_room_id ON sensors(room_id);

-- Index to speed up measurement lookup by sensor and time
CREATE INDEX idx_measurements_sensor_time ON measurements(sensor_id, timestamp);
```

### Insert statements
```sql
INSERT INTO rooms (name) VALUES ('room_A'), ('room_B');

-- room_A sensors
INSERT INTO sensors (name, room_id, parameter_type)
VALUES
  ('sensor_A_v1', (SELECT id FROM rooms WHERE name = 'room_A'), 'V'),
  ('sensor_A_r1', (SELECT id FROM rooms WHERE name = 'room_A'), 'R'),
  ('sensor_A_r2', (SELECT id FROM rooms WHERE name = 'room_A'), 'R');


-- room_B sensors
INSERT INTO sensors (name, room_id, parameter_type)
VALUES
  ('sensor_B_v1', (SELECT id FROM rooms WHERE name = 'room_B'), 'V'),
  ('sensor_B_v2', (SELECT id FROM rooms WHERE name = 'room_B'), 'V'),
  ('sensor_B_r1', (SELECT id FROM rooms WHERE name = 'room_B'), 'R'),
  ('sensor_B_r2', (SELECT id FROM rooms WHERE name = 'room_B'), 'R'),
  ('sensor_B_r3', (SELECT id FROM rooms WHERE name = 'room_B'), 'R');

-- Assume now() is 2025-06-27 20:00:00+00
-- V and R measurements in room_A
INSERT INTO measurements (sensor_id, timestamp, value)
VALUES
  -- room_A V sensor sends data
  ((SELECT id FROM sensors WHERE name = 'sensor_A_v1'), '2025-06-27 20:00:00+00', 220.0),

  -- only one R sensor in room_A sends data
  ((SELECT id FROM sensors WHERE name = 'sensor_A_r1'), '2025-06-27 20:00:00+00', 10.5);

-- Only R data in room_B
INSERT INTO measurements (sensor_id, timestamp, value)
VALUES
  ((SELECT id FROM sensors WHERE name = 'sensor_B_r1'), '2025-06-27 20:00:00+00', 12.1),
  ((SELECT id FROM sensors WHERE name = 'sensor_B_r2'), '2025-06-27 20:00:00+00', 12.3);

-- Only V data later in room_B
INSERT INTO measurements (sensor_id, timestamp, value)
VALUES
  ((SELECT id FROM sensors WHERE name = 'sensor_B_v1'), '2025-06-27 20:01:00+00', 230.0),
  ((SELECT id FROM sensors WHERE name = 'sensor_B_v2'), '2025-06-27 20:01:00+00', 229.5);

```

### Select statement
```sql
WITH measurements_with_info AS (
    SELECT
        r.name AS room,
        date_trunc('second', m.timestamp) AS ts,
        s.parameter_type,
        m.value
    FROM measurements m
    JOIN sensors s ON m.sensor_id = s.id
    JOIN rooms r ON s.room_id = r.id
),
aggregated AS (
    SELECT
        room,
        ts AS timestamp,
        AVG(CASE WHEN parameter_type = 'V' THEN value END) AS avg_v,
        AVG(CASE WHEN parameter_type = 'R' THEN value END) AS avg_r
    FROM measurements_with_info
    GROUP BY room, ts
)
SELECT
    room,
    timestamp,
    avg_v / NULLIF(avg_r, 0) AS I,
    avg_v AS V,
    avg_r AS R
FROM aggregated
ORDER BY room, timestamp;
```