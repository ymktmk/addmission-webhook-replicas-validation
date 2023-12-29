
## Webhookの実装
Kubernetesでは、リソースの作成・更新・削除をおこなう直前にWebhookで任意の処理を実行するとことができます。MutatingWebhookではリソースの値を書き換えることができ、ValidatingWebhookでは値の検証をおこなうことができます。

Ref: https://zoetrope.github.io/kubebuilder-training/controller-runtime/webhook.html

## オレオレ証明書

### webhook serverの秘密鍵作成
```
openssl genrsa 2048 > server.key
```

### CSRの作成
```
openssl req -new -key server.key -out server.csr
```

### SubjectAltNameに対応するためのファイル作成
````
echo "subjectAltName = DNS:replicas-validating-webhook.default.svc, DNS:replicas-validating-webhook.default.svc.cluster.local" > san.txt
````

### サーバ証明書の作成
```
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt -extfile san.txt
```

## ローカルで起動する

```
go run main.go -server-cert=./server.crt -server-key=./server.key
```

### Secret 
```
kubectl create secret tls replicas-validating-webhook-secret --key server.key --cert server.crt
```

### デプロイ
```
sed  "s/BASE64_ENCODED_PEM_FILE/$(base64 -i server.crt)/g" manifests/validatingwebhookconfiguration.yaml.template | kubectl apply -f -
```
