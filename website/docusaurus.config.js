/**
 * @type {import('@docusaurus/types').DocusaurusConfig}
 */
module.exports = {
  baseUrl: '/',
  favicon: 'img/manael.png',
  i18n: {
    defaultLocale: 'en',
    localeConfigs: {
      en: {
        label: 'English'
      },
      ja: {
        label: '日本語'
      }
    },
    locales: ['en', 'ja']
  },
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  organizationName: 'manaelproxy',
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js')
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css')
        }
      }
    ]
  ],
  projectName: 'manael',
  tagline: 'Manael is a simple HTTP proxy for processing images.',
  themeConfig: {
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
  },
  title: 'Manael',
  trailingSlash: false,
  url: 'https://manael.org'
}
