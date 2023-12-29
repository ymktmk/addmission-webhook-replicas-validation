
## Webhookの実装
Kubernetesでは、リソースの作成・更新・削除をおこなう直前にWebhookで任意の処理を実行するとことができます。MutatingWebhookではリソースの値を書き換えることができ、ValidatingWebhookでは値の検証をおこなうことができます。

Ref: https://zoetrope.github.io/kubebuilder-training/controller-runtime/webhook.html
