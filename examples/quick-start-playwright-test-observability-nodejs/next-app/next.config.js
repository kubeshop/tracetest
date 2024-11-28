/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    instrumentationHook: true,
  },
  pageExtensions: ['mdx', 'md', 'jsx', 'js', 'tsx', 'ts', 'node.ts'],
  output: 'standalone',
}

module.exports = nextConfig
