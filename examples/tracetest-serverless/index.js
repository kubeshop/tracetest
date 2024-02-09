const https = require('https');

function getRequest() {
  const url = 'https://tracetest.io/';

  return new Promise((resolve, reject) => {
    const req = https.get(url, (res) => {
      resolve(res.statusCode);
    });

    req.on('error', (err) => {
      reject(new Error(err));
    });
  });
}

exports.handler = async (event) => {
  try {
    const result = await getRequest();
    return {
      statusCode: result,
    };
  } catch (error) {
    return {
      statusCode: 400,
      body: error.message,
    };
  }
};
