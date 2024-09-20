module.exports.hello = async (event) => {
  const response = {
    statusCode: 200,
    body: JSON.stringify({
      message: 'Hello, world!',
    }),
  };

  return response;
};
