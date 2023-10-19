package quick_start_api;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.Banner;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;

import jakarta.servlet.Filter;

import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.sdk.OpenTelemetrySdk;
import io.opentelemetry.sdk.autoconfigure.AutoConfiguredOpenTelemetrySdk;
import io.opentelemetry.instrumentation.spring.webmvc.v6_0.SpringWebMvcTelemetry;

@SpringBootApplication
public class QuickStartApplication {
  private static volatile OpenTelemetry openTelemetry = OpenTelemetry.noop();

  public static void main(String[] args) {
    OpenTelemetrySdk openTelemetrySdk = AutoConfiguredOpenTelemetrySdk.builder().build().getOpenTelemetrySdk();
    QuickStartApplication.openTelemetry = openTelemetrySdk;

    SpringApplication app = new SpringApplication(QuickStartApplication.class);
    app.setBannerMode(Banner.Mode.OFF);
    app.run(args);
  }

  @Bean
  public OpenTelemetry openTelemetry() {
    return openTelemetry;
  }

  @Bean
  public Filter webMvcTracingFilter(OpenTelemetry openTelemetry) {
      return SpringWebMvcTelemetry.create(openTelemetry).createServletFilter();
  }
}
