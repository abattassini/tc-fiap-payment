terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
  
  # AWS Academy usa credenciais temporárias via environment variables
  # Não configurar access_key/secret_key aqui
}

# RDS PostgreSQL - Payment Database
resource "aws_db_instance" "payment" {
  identifier           = var.db_instance_identifier
  engine               = "postgres"
  engine_version       = "15.4"
  instance_class       = var.db_instance_class
  allocated_storage    = 20
  max_allocated_storage = 100  # Auto-scaling até 100GB
  storage_type         = "gp3"
  storage_encrypted    = true
  
  # Database configuration
  db_name  = var.db_name
  username = var.db_username
  password = var.db_password
  port     = 5432
  
  # Network configuration
  publicly_accessible = true  # Para AWS Academy (ajustar em produção)
  skip_final_snapshot = true  # Para facilitar testes (remover em produção)
  
  # Backup configuration
  backup_retention_period = 7
  backup_window          = "03:00-04:00"
  maintenance_window     = "mon:04:00-mon:05:00"
  
  # Performance Insights (pode não estar disponível no Academy)
  # performance_insights_enabled = true
  # performance_insights_retention_period = 7
  
  tags = {
    Name        = "Payment Database"
    Environment = var.environment
    Project     = "tc-fiap-payment"
    ManagedBy   = "Terraform"
  }
}

# Output útil para o pipeline
output "rds_endpoint" {
  description = "Endpoint da instância RDS PostgreSQL"
  value       = aws_db_instance.payment.endpoint
  sensitive   = true
}

output "rds_address" {
  description = "Endereço da instância RDS PostgreSQL (sem porta)"
  value       = aws_db_instance.payment.address
}

output "rds_port" {
  description = "Porta da instância RDS PostgreSQL"
  value       = aws_db_instance.payment.port
}

output "rds_db_name" {
  description = "Nome do banco de dados"
  value       = aws_db_instance.payment.db_name
}

output "rds_arn" {
  description = "ARN da instância RDS"
  value       = aws_db_instance.payment.arn
}
