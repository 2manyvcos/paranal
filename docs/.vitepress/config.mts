import { defineConfig } from 'vitepress';

// https://vitepress.dev/reference/site-config
export default defineConfig({
  vite: {
    server: { host: '0.0.0.0', port: 8080, allowedHosts: true },
  },

  lang: 'en-US',
  title: 'Paranal',
  description: 'Advanced Service Dashboard with Health and Version Monitoring',

  lastUpdated: true,
  cleanUrls: true,

  head: [
    ['link', { rel: 'icon', href: '/paranal-simple.svg' }],
    ['link', { rel: 'apple-touch-icon', href: '/apple-touch-icon.png' }],
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:site_name', content: 'Paranal' }],
    [
      'meta',
      {
        property: 'og:title',
        content:
          'Advanced Service Dashboard with Health and Version Monitoring',
      },
    ],
    [
      'meta',
      { property: 'og:url', content: 'https://2manyvcos.github.io/paranal' },
    ],
    // ["meta", { property: "og:image", content: "https://2manyvcos.github.io/paranal/social.png" }],
  ],

  themeConfig: {
    logo: '/paranal-simple.svg',

    search: {
      provider: 'local',
    },

    editLink: {
      pattern: 'https://github.com/2manyvcos/paranal/edit/main/docs/:path',
    },

    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Guide', link: '/guide/what-is-paranal', activeMatch: '/guide/' },
      { text: 'API Reference', link: '/api' },
    ],

    sidebar: {
      '/guide/': [
        {
          text: 'Introduction',
          items: [
            { text: 'What is Paranal?', link: '/guide/what-is-paranal' },
            { text: 'Getting Started', link: '/guide/getting-started' },
          ],
        },
      ],
    },

    footer: {
      message: 'Released under the MIT License',
      copyright:
        "Copyright © 2026-present <a href='https://github.com/2manyvcos'>Aaron Burmeister</a>",
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/2manyvcos/paranal' },
      { icon: 'docker', link: 'https://hub.docker.com/r/2manyvcos/paranal' },
    ],
  },

  sitemap: {
    hostname: 'https://2manyvcos.github.io/paranal',
  },

  base: '/paranal/',
});
