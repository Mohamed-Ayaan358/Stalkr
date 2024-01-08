/** @type {import('next').NextConfig} */
module.exports = {
  async rewrites() {
    return [
      {
        source: '/api/add',
        destination: 'http://localhost:8080/add',
      },
    ];
  },
};
