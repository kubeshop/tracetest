const handler = async (event) => {
  console.log("Event: ", event);
  let responseMessage = "Hello, World!";

  return {
    statusCode: 200,
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      message: responseMessage,
    }),
  };
};

module.exports.handler = handler;
