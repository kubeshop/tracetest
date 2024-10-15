package com.example;

import static spark.Spark.*;

public class HelloWorldApi {
    public static void main(String[] args) {
        // Start a basic HTTP server on port 4567
        port(4567);

        // Define a simple route that returns "Hello, World!"
        get("/hello", (req, res) -> "Hello, World!");
    }
}
