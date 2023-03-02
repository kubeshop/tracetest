output "api_endpoint" {
  value       = "${aws_apigatewayv2_stage.lambda.invoke_url}/hello"
  description = "The API endpoint"
}

output "tracetest_url" {
  value       = "http://${aws_lb.tracetest-alb.dns_name}:11633"
  description = "Tracetest public URL"
}

output "jaeger_ui_url" {
  value       = "http://${aws_lb.tracetest-alb.dns_name}:16686"
  description = "Jaeger public URL"
}

output "internal_jaeger_api_url" {
  value       = "${aws_lb.internal_tracetest_alb.dns_name}:16685"
  description = "Jaeger internal API URL"
}
