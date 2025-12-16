# ESSP Monitoring Stack

This directory contains the monitoring configuration for the ESSP microservices platform using Prometheus, Grafana, and AlertManager.

## Table of Contents

- [Overview](#overview)
- [Components](#components)
- [Quick Start](#quick-start)
- [Available Metrics](#available-metrics)
- [Dashboards](#dashboards)
- [Alert Rules](#alert-rules)
- [Kubernetes Deployment](#kubernetes-deployment)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)

## Overview

The ESSP monitoring stack provides:

- **Metrics Collection**: Prometheus scrapes metrics from all ESSP services
- **Visualization**: Grafana dashboards for observability
- **Alerting**: Prometheus alerts with AlertManager for notifications
- **Recording Rules**: Pre-computed metrics for faster queries

## Components

### Prometheus

Prometheus is the metrics collection and storage system. It scrapes metrics from:

- IMS API service
- SSOT services (devices, parts, school)
- Sync Worker service
- PostgreSQL database (via postgres-exporter)
- Redis cache (via redis-exporter)
- NATS messaging (native metrics)

**Port**: 9090
**Config**: `prometheus/prometheus.yml`

### Grafana

Grafana provides visualization and dashboards for metrics.

**Port**: 3000
**Default Credentials**: admin/admin
**Dashboards**: `grafana/dashboards/`

### AlertManager

AlertManager handles alerts sent by Prometheus and routes them to various notification channels.

**Port**: 9093
**Config**: `alertmanager/config.yml`

## Quick Start

### Local Development

1. **Start the monitoring stack**:

   ```bash
   cd deployments
   docker-compose -f docker-compose.monitoring.yml up -d
   ```

2. **Access the services**:
   - Prometheus: http://localhost:9090
   - Grafana: http://localhost:3000 (admin/admin)
   - AlertManager: http://localhost:9093

3. **Configure AlertManager** (optional):

   Edit `monitoring/alertmanager/config.yml` to add your notification channels (Slack, email, PagerDuty, etc.)

4. **Stop the monitoring stack**:

   ```bash
   docker-compose -f docker-compose.monitoring.yml down
   ```

### Kubernetes Deployment

1. **Install Prometheus Operator** (if not already installed):

   ```bash
   kubectl create -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/main/bundle.yaml
   ```

2. **Deploy ServiceMonitors**:

   ```bash
   kubectl apply -f monitoring/servicemonitor.yaml
   ```

3. **Verify ServiceMonitors**:

   ```bash
   kubectl get servicemonitors -n essp
   ```

## Available Metrics

### HTTP Metrics

All ESSP services expose these standard HTTP metrics:

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `http_requests_total` | Counter | Total HTTP requests | method, path, status |
| `http_request_duration_seconds` | Histogram | HTTP request duration | method, path |

### Database Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `db_connections_active` | Gauge | Active database connections |

### Business Metrics

IMS API specific business metrics:

| Metric | Type | Description |
|--------|------|-------------|
| `incidents_created_total` | Counter | Total incidents created |
| `work_orders_created_total` | Counter | Total work orders created |

### Recording Rules

Pre-computed metrics for performance:

| Metric | Description |
|--------|-------------|
| `essp:http_requests:rate5m` | Request rate per service (5m) |
| `essp:http_error_ratio:rate5m` | Error ratio per service (5m) |
| `essp:http_request_duration:p95` | P95 latency per service |
| `essp:http_request_duration:p99` | P99 latency per service |
| `essp:platform_availability:ratio5m` | Overall platform availability |

See `prometheus/recording-rules.yml` for complete list.

## Dashboards

### ESSP Platform Overview

**File**: `grafana/dashboards/essp-overview.json`
**UID**: `essp-overview`

Main dashboard showing:
- Platform availability
- Request rate by service
- Platform error rate
- Request latency (p95, p99)
- Error rate by service
- Business metrics (incidents, work orders)
- Database connections

### IMS API Metrics

**File**: `grafana/dashboards/ims-api.json`
**UID**: `ims-api`

IMS API specific dashboard showing:
- Request rate by HTTP method
- Response status codes distribution
- Request latency percentiles (p50, p95, p99)
- Top 10 slowest endpoints
- Top 10 endpoints by request rate
- Business metrics (incidents, work orders)

### Importing Dashboards

Dashboards are automatically provisioned when using docker-compose. For manual import:

1. Navigate to Grafana (http://localhost:3000)
2. Click + → Import
3. Upload the JSON file from `grafana/dashboards/`

## Alert Rules

### Critical Alerts

| Alert | Condition | Duration | Description |
|-------|-----------|----------|-------------|
| `HighErrorRate` | 5xx errors > 5% | 5m | High server error rate |
| `VeryHighLatency` | p95 > 1s | 3m | Very high request latency |
| `ServiceDown` | Service unavailable | 1m | Service is down |
| `DBConnectionPoolCritical` | Pool > 95% full | 2m | Database connection exhaustion |
| `PostgreSQLDown` | PostgreSQL unavailable | 1m | Database is down |
| `CriticalMemoryUsage` | Memory > 95% | 2m | Critical memory usage |

### Warning Alerts

| Alert | Condition | Duration | Description |
|-------|-----------|----------|-------------|
| `HighLatency` | p95 > 500ms | 5m | High request latency |
| `DBConnectionPoolExhaustion` | Pool > 80% full | 5m | Database pool near exhaustion |
| `HighMemoryUsage` | Memory > 85% | 5m | High memory usage |
| `HighCPUUsage` | CPU > 80% | 10m | High CPU usage |

### Business Alerts

| Alert | Condition | Duration | Description |
|-------|-----------|----------|-------------|
| `IncidentCreationAnomaly` | Rate > 2x avg | 15m | Unusual incident creation rate |
| `WorkOrderCreationAnomaly` | Rate > 2x avg | 15m | Unusual work order creation |

See `prometheus/alerts.yml` for complete alert definitions.

## Kubernetes Deployment

### Prerequisites

- Kubernetes cluster with Prometheus Operator installed
- Services labeled with `app.kubernetes.io/part-of: essp`

### ServiceMonitor Configuration

The ServiceMonitor resources automatically discover and scrape ESSP services:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: essp-services
  labels:
    release: prometheus
spec:
  selector:
    matchLabels:
      app.kubernetes.io/part-of: essp
  endpoints:
  - port: http
    path: /metrics
    interval: 15s
```

### Viewing Metrics in Production

```bash
# Port-forward Prometheus
kubectl port-forward -n monitoring svc/prometheus-operated 9090:9090

# Port-forward Grafana
kubectl port-forward -n monitoring svc/grafana 3000:3000
```

## Configuration

### Adding New Metrics

1. **Define the metric** in `services/ims-api/internal/metrics/metrics.go`:

   ```go
   var MyNewMetric = prometheus.NewCounter(prometheus.CounterOpts{
       Name: "my_new_metric_total",
       Help: "Description of my new metric",
   })
   ```

2. **Register the metric** in the `init()` function:

   ```go
   prometheus.MustRegister(MyNewMetric)
   ```

3. **Increment the metric** in your code:

   ```go
   metrics.MyNewMetric.Inc()
   ```

### Adding New Alert Rules

1. Edit `prometheus/alerts.yml`
2. Add your alert under the appropriate group
3. Restart Prometheus to reload configuration:

   ```bash
   docker-compose -f docker-compose.monitoring.yml restart prometheus
   ```

### Configuring Notifications

Edit `alertmanager/config.yml` to add notification receivers:

**Slack Example**:
```yaml
receivers:
  - name: 'critical-receiver'
    slack_configs:
      - api_url: 'YOUR_SLACK_WEBHOOK_URL'
        channel: '#essp-critical-alerts'
        title: 'Critical Alert: {{ .GroupLabels.alertname }}'
```

**Email Example**:
```yaml
receivers:
  - name: 'critical-receiver'
    email_configs:
      - to: 'oncall@essp.example.com'
        headers:
          Subject: 'Critical Alert: {{ .GroupLabels.alertname }}'
```

**PagerDuty Example**:
```yaml
receivers:
  - name: 'critical-receiver'
    pagerduty_configs:
      - service_key: 'YOUR_PAGERDUTY_SERVICE_KEY'
```

## Troubleshooting

### Metrics Not Appearing

1. **Check service is exposing /metrics endpoint**:
   ```bash
   curl http://localhost:8080/metrics
   ```

2. **Check Prometheus targets**:
   - Navigate to http://localhost:9090/targets
   - Verify all targets are "UP"

3. **Check Prometheus logs**:
   ```bash
   docker-compose -f docker-compose.monitoring.yml logs prometheus
   ```

### Dashboards Not Loading

1. **Verify Prometheus datasource**:
   - Navigate to Grafana → Configuration → Data Sources
   - Test the Prometheus connection

2. **Check Grafana logs**:
   ```bash
   docker-compose -f docker-compose.monitoring.yml logs grafana
   ```

### Alerts Not Firing

1. **Check alert rules are loaded**:
   - Navigate to http://localhost:9090/rules
   - Verify your alert rules are listed

2. **Check AlertManager**:
   - Navigate to http://localhost:9093
   - Verify alerts are being received

3. **Test alert expression**:
   - Navigate to http://localhost:9090/graph
   - Run the alert expression manually

### High Cardinality Issues

If you experience performance issues due to high cardinality metrics:

1. **Reduce label diversity**: Avoid dynamic labels with many unique values
2. **Use recording rules**: Pre-compute common queries
3. **Adjust retention**: Reduce Prometheus retention period
4. **Use relabeling**: Drop unnecessary labels in Prometheus scrape config

## Best Practices

1. **Metric Naming**: Follow Prometheus naming conventions
   - Use `_total` suffix for counters
   - Use `_seconds` for durations
   - Use base units (seconds, bytes)

2. **Labels**: Keep label cardinality low
   - Don't use user IDs or request IDs as labels
   - Use labels for dimensions you'll query by

3. **Dashboards**:
   - Use templating for multi-service dashboards
   - Set appropriate refresh intervals (30s recommended)
   - Include documentation panels

4. **Alerts**:
   - Use appropriate thresholds and durations
   - Include actionable information in annotations
   - Set up proper alert routing and inhibition

## Additional Resources

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [AlertManager Documentation](https://prometheus.io/docs/alerting/latest/alertmanager/)
- [Prometheus Operator](https://prometheus-operator.dev/)
- [Best Practices for Monitoring](https://prometheus.io/docs/practices/)

## Support

For questions or issues with the monitoring stack, contact the platform team or create an issue in the repository.
