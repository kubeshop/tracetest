output "tracetest_url" {
  value       = "http://${aws_lb.tracetest_alb.dns_name}:11633"
  description = "Tracetest public URL"
}
