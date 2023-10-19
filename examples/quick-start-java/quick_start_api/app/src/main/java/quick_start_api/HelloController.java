package quick_start_api;

import java.util.List;
import java.util.Optional;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.server.ResponseStatusException;

import org.springframework.beans.factory.annotation.Autowired;
import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.api.trace.Span;
import io.opentelemetry.context.Scope;

@RestController
public class HelloController {
  private final Tracer tracer;

  @GetMapping("/hello")
  public String index() {
    Span span = tracer.spanBuilder("hello").startSpan();

    try (Scope scope = span.makeCurrent()) {
      return "hello";
    } catch(Throwable t) {
      span.recordException(t);
      throw t;
    } finally {
      span.end();
    }
  }

  @Autowired
  HelloController(OpenTelemetry openTelemetry) {
    tracer = openTelemetry.getTracer(HelloController.class.getName(), "0.1.0");
  }
}
