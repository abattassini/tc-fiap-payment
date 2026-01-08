variable "aws_region" {
  description = "AWS Region"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name (staging, production, etc)"
  type        = string
  default     = "staging"
}

# RDS Configuration
variable "db_instance_identifier" {
  description = "Identificador da instância RDS"
  type        = string
  default     = "tc-fiap-payment-db"
}

variable "db_instance_class" {
  description = "Classe da instância RDS (AWS Academy limita a db.t3.micro)"
  type        = string
  default     = "db.t3.micro"
}

variable "db_name" {
  description = "Nome do banco de dados PostgreSQL"
  type        = string
  default     = "payment_db"
}

variable "db_username" {
  description = "Nome de usuário do banco de dados"
  type        = string
  default     = "payment_user"
  sensitive   = true
}

variable "db_password" {
  description = "Senha do banco de dados"
  type        = string
  sensitive   = true
}
