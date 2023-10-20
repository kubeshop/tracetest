package risk.analysis.worker;

import java.time.Duration;
import java.util.Arrays;
import java.util.Properties;
import org.apache.kafka.clients.consumer.Consumer;
import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.apache.kafka.common.serialization.StringDeserializer;

import com.fasterxml.jackson.databind.ObjectMapper;

import io.opentelemetry.api.GlobalOpenTelemetry;
import io.opentelemetry.api.trace.Span;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.context.Scope;

public class App {
    public static void main(String[] args) {
        String bootstrapServers = System.getenv("KAFKA_BROKER_URL");
        String groupId = System.getenv("OTEL_SERVICE_NAME");
        String topic = System.getenv("KAFKA_TOPIC");

        // create consumer configs
        Properties properties = new Properties();
        properties.setProperty(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, bootstrapServers);
        properties.setProperty(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        properties.setProperty(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        properties.setProperty(ConsumerConfig.GROUP_ID_CONFIG, groupId);
        properties.setProperty(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");

        System.out.println("Setting up worker...");

        System.out.println("Initializing Kafka Consumer...");
        Consumer<String, String> consumer = new KafkaConsumer<>(properties);
        consumer.subscribe(Arrays.asList(topic));
        System.out.println("Kafka Consumer initialized.");

        System.out.println("Polling for records on Kafka stream...");

        while (true){
          ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(100));

          for (ConsumerRecord<String, String> record : records){
            System.out.println("Receiving message from Partition: " + record.partition() + ", Offset:" + record.offset());

            String paymentOrderAsJSON = record.value();
            analyseOrder(paymentOrderAsJSON);
          }
      }
    }

    private static void analyseOrder(String paymentOrderAsJSON) {
      Tracer tracer = GlobalOpenTelemetry.getTracer("analyze-order-tracer");

      Span span = tracer.spanBuilder("analyseOrder").startSpan();

      try (Scope scope = span.makeCurrent()) {
        ObjectMapper mapper = new ObjectMapper();
        PaymentOrder order = mapper.readValue(paymentOrderAsJSON, PaymentOrder.class);

        Boolean isHighRiskRate = (order.getValue() > 10000.0);

        span.setAttribute("riskAnalysis.highRiskRate", isHighRiskRate);

        if (isHighRiskRate) {
          System.out.println("Order from customer " + order.getOriginCustomerID() + " to customer " + order.getDestinationCustomerID() + " has a high risk rate");
        } else {
          System.out.println("Order from customer " + order.getOriginCustomerID() + " to customer " + order.getDestinationCustomerID() + " is ok");
        }
      } catch (Exception e) {
        span.recordException(e);
        System.out.println("Error when reading paymentOrder: " + e.toString());
      } finally {
        span.end();
      }
    }
}

