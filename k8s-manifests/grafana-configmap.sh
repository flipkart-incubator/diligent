kubectl create configmap grafana-prov-datasources --from-file=../monitoring/grafana/provisioning/datasources
kubectl create configmap grafana-prov-dashboards --from-file=../monitoring/grafana/provisioning/dashboards
kubectl create configmap grafana-diligent-dashboards --from-file=../monitoring/grafana/dashboards/diligent
