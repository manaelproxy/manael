module.exports = {
  baseUrl: '/',
  favicon: 'img/manael.png',
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
  url: 'https://manael.org'
}
