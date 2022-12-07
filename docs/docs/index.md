# Tracetest Documentation

Generate end-to-end tests automatically from your traces. For QA, Dev, & Ops.

<!-- 
TODO: migrate video to youtube and use YT embed.

<p align="center">
 <script src="https://fast.wistia.com/embed/medias/dw06408oqz.jsonp" async></script><script src="https://fast.wistia.com/assets/external/E-v1.js" async></script><div class="wistia_responsive_padding" style="padding:56.25% 0 0 0;position:relative;"><div class="wistia_responsive_wrapper" style="height:100%;left:0;position:absolute;top:0;width:100%;"><div class="wistia_embed wistia_async_dw06408oqz videoFoam=true" style="height:100%;position:relative;width:100%"><div class="wistia_swatch" style="height:100%;left:0;opacity:0;overflow:hidden;position:absolute;top:0;transition:opacity 200ms;width:100%;"><img src="https://fast.wistia.com/embed/medias/dw06408oqz/swatch" style="filter:blur(5px);height:100%;object-fit:contain;width:100%;" alt="" aria-hidden="true" onload="this.parentNode.style.opacity=1;" /></div></div></div></div>
</p>

-->

Tracetest allows you to quickly build integration and end-to-end tests, powered by your OpenTelemetry traces.

- Point Tracetest to your preffered trace back-end, like Jaeger or Tempo, or to the OpenTelemetry Collector directly.
- Define a triggering transaction, such as a GET against an API endpoint.
- The system runs this transaction, returning both the response data and a full trace.
- Define tests & assertions against this data, ensuring both your response and the underlying processes worked correctly, quickly, and without errors.
- Save your test.
- Run the tests either manually or via your CI build jobs.

New to trace-based testing? Read more about the concepts, [here](./concepts/introduction-to-trace-based-testing).
