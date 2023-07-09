# webhook-bb-pr-teams-router-app-go


To run this app and test locally follow this:

```bash
cd backend-test
./prepare-ca.sh
./localbuild.sh
docker run -it --rm -p 80:8080 -e RLOG_LOG_LEVEL=DEBUG backend-test
docker run -it --rm -e "HOST=$(hostname -I | cut -d' ' -f1)" -p 443:443 tls-test
cd adaptor
./localbuild.sh

docker run -it --rm -e "TEAMS_HOSTNAME=$(hostname -I | cut -d' ' -f1)" -e RLOG_LOG_LEVEL=INFO -e TLS_INSECURE_SKIP_VERIFY=false -p 8080:8080 docker.io/library/adaptor

curl -v -X POST -H 'Content-Type: application/json' -H "X-Request-Id: $(uuidgen)" "http://127.0.0.1:8080/webhookb2/$(uuidgen)@$(uuidgen)/IncomingWebhook/$(openssl rand -hex 16)/$(uuidgen)" -d @bb-hook-event.json

```