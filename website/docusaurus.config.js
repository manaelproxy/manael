// @ts-check
// `@type` JSDoc annotations allow editor autocompletion and type checking
// (when paired with `@ts-check`).
// There are various equivalent ways to declare your Docusaurus config.
// See: https://docusaurus.io/docs/api/docusaurus-config

/** @type {import('@docusaurus/types').Config} */
const config = {
  baseUrl: '/',
  favicon: 'img/manael.png',
  i18n: {
    defaultLocale: 'en',
    locales: ['en', 'ja']
  },
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  organizationName: 'manaelproxy',
  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: './sidebars.js'
        },
        theme: {
          customCss: './src/css/custom.css'
        }
      })
    ]
  ],
  projectName: 'manael',
  tagline: 'Manael is a simple HTTP proxy for processing images.',
  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      algolia: {
        appId: 'ZX7VYOHRJ3',
        apiKey: '43f66f766ffb77ee2280608d793ab235',
        indexName: 'docusaurus'
      },
      footer: {
        copyright: 'Copyright © 2018 The Manael Authors.',
        links: [],
        style: 'dark'
      },
      navbar: {
        hideOnScroll: true,
        items: [
          {
            activeBasePath: 'docs',
            label: 'Docs',
            position: 'left',
            to: 'docs/'
          },
          {
            type: 'localeDropdown',
            position: 'right'
          },
          {
            href: 'https://github.com/manaelproxy/manael',
            label: 'GitHub',
            position: 'right'
          }
        ],
        logo: {
          alt: 'Manael Logo',
          src: 'img/manael.png'
        },
        title: 'Manael'
      }
    }),
  title: 'Manael',
  trailingSlash: false,
  url: 'https://manael.org'
}

export default config
