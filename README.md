# teleprobe
teleprobe is a lightweight telemetry system that simulates sensor nodes sending real-time data to a telemetry sink

## Architecture Overview

+------------------+     gRPC / HTTP / MQTT / NATS     +--------------------+
|  Sensor Node     | --------------------------------> |  Telemetry Sink    |
|  (Data Emitter)  |                                   |  (Data Receiver)   |
+------------------+                                   +--------------------+
        ↑                                                      ↓
        |                                              DB / File / Message Queue
        | (Retry, Buffer)                                   (Storage, Analytics)
        |
    Local Buffer

