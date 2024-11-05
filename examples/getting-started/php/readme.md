# Step-by-step

1. Install OpenTelemetry Extension for PHP

      ```bash
      pecl install opentelemetry
      ```

2. Add OpenTelemetry to your `php.ini`

      ```bash
      [opentelemetry]
      extension=opentelemetry.so
      ```

3. Install `composer` if not installed already

      ```bash
      php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');"
      php composer-setup.php
      php -r "unlink('composer-setup.php');"
      
      # If you want to use Composer globally from any directory, move it to a directory in your PATH:
      sudo mv composer.phar /usr/local/bin/composer
      composer --version
      ```

4. Install OpenTelemetry SDK for your libraries

      [View all libraries here.](https://packagist.org/search/?query=open-telemetry&tags=instrumentation) The example below showcases Laravel.

      ```bash
      composer require \
        open-telemetry/sdk \
        open-telemetry/exporter-otlp \
        open-telemetry/opentelemetry-auto-laravel
        # ...        
      ```

5. Configure OpenTelemetry env vars in `php.ini`

      ```ini
      [opentelemetry]
      extension=opentelemetry.so

      OTEL_PHP_AUTOLOAD_ENABLED="true"
      OTEL_SERVICE_NAME=your-php-app
      OTEL_TRACES_EXPORTER=otlp
      OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
      OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
      OTEL_PROPAGATORS=baggage,tracecontext
      ```
