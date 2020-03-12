module.exports = {
  baseUrl: '/',
  favicon: 'img/manael.png',
  organizationName: 'manaelproxy',
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          path: '../docs',
          showLastUpdateAuthor: true,
          showLastUpdateTime: true,
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
      links: [
        {
          items: [
            {
              label: 'Introduction',
              to: 'docs/introduction'
            },
            {
              label: 'Installation',
              to: 'docs/installation'
            }
          ],
          title: 'Docs'
        }
      ],
      style: 'dark'
    },
    navbar: {
      hideOnScroll: true,
      links: [
        {
          activeBasePath: 'docs',
          label: 'Docs',
          position: 'left',
          to: 'docs/introduction'
        },
        {
          href: 'https://github.com/manaelproxy/manael',
          label: 'GitHub',
          position: 'right'
        }
      ],
      logo: {
        alt: '',
        src: 'img/manael.png'
      },
      title: 'Manael'
    }
  },
  title: 'Manael',
  url: 'https://manael.org'
}
