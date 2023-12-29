
## Webhookの実装
Kubernetesでは、リソースの作成・更新・削除をおこなう直前にWebhookで任意の処理を実行するとことができます。MutatingWebhookではリソースの値を書き換えることができ、ValidatingWebhookでは値の検証をおこなうことができます。

Ref: https://zoetrope.github.io/kubebuilder-training/controller-runtime/webhook.html


kubectl create secret tls replicas-validating-webhook-secret --key server.key --cert server.crt
kubectl describe secret replicas-validating-webhook-secret

## AdmissionWebhook

sed  "s/BASE64_ENCODED_PEM_FILE/$(base64 server.crt)/g" manifests/validatingwebhookconfiguration.yaml.template | kubectl apply -f -
