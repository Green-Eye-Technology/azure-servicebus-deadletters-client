# azure-servicebus-deadletters-client

A Go client that handles Azure service bus dead letters

## Usage

```
SERVICEBUS_CONNECTION_STRING="Endpoint=https://<servicebus-name>.servicebus.windows.net/;SharedAccessKeyName=YOUR_SHARED_ACCESS_KEY_NAME;SharedAccessKey=YOUR_SHARE_ACCESS_KEY" go run main.go --subscriptionName <subscription-name> --topicName <topic-name>
```
