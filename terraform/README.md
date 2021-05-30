# Terraform

This module sets up Kafka and PostgreSQL in [Aiven](https://aiven.io/). If you don't have an account you can [sign up](https://console.aiven.io/signup) for 30 day trial and get free $300 credits. Credit card is not required for the trial.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.14+ (earlier might work too, but untested)
- Aiven API token (generate one at https://console.aiven.io/profile/auth)
- Aiven project name

## Usage

```bash
export TF_VAR_aiven_project_name="YOUR AIVEN PROJECT NAME"
export TF_VAR_aiven_api_token="YOUR_AIVEN_API_TOKEN"

terraform init
terraform apply
```

We need some secrets for the Kafka and PostgreSQL:

```bash
# Kafka
terraform output -raw kafka_access_key > kafka_access.key
terraform output -raw kafka_access_cert > kafka_access.cert
terraform output -raw kafka_service_uri > kafka_service_uri
terraform output -raw project_ca_cert > ca.pem

# PostgreSQL
terraform output -raw postgresql_service_uri > postgres_dburl
```

## Links

- [Terraform Aiven provider](https://github.com/aiven/terraform-provider-aiven) check examples folder for code samples
