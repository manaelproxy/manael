import PropTypes from 'prop-types'
import React from 'react'
import classnames from 'classnames'
import Layout from '@theme/Layout'
import Link from '@docusaurus/Link'
import useDocusaurusContext from '@docusaurus/useDocusaurusContext'
import useBaseUrl from '@docusaurus/useBaseUrl'
import styles from './styles.module.css'

const features = [
  {
    description: <>Just run a one binary!</>,
    title: <>Simple!</>
  },
  {
    description: <>Manael’s binary run anywhere in a GNU/Linux environment.</>,
    title: <>Portability!</>
  },
  {
    description: (
      <>Manael is fast because there is no unnecessary processing!</>
    ),
    title: <>High Peformance!</>
  }
]

function Feature({ description, imageUrl, title }) {
  const imgUrl = useBaseUrl(imageUrl)
  return (
    <div className={classnames('col col--4', styles.feature)}>
      {imgUrl && (
        <div className="text--center">
          <img className={styles.featureImage} src={imgUrl} alt={title} />
        </div>
      )}
      <h3>{title}</h3>
      <p>{description}</p>
    </div>
  )
}

Feature.propTypes = {
  description: PropTypes.string.isRequired,
  imageUrl: PropTypes.string,
  title: PropTypes.string.isRequired
}

function Home() {
  const context = useDocusaurusContext()
  const { siteConfig = {} } = context

  return (
    <Layout description="Manael is a simple HTTP proxy for processing images.">
      <header className={classnames('hero hero--dark', styles.heroBanner)}>
        <div className="container">
          <h1 className="hero__title">
            <img
              alt={siteConfig.title}
              className={classnames('margin-vert--md', styles.heroBannerLogo)}
              height={128}
              src={useBaseUrl('img/logo.png')}
              width={128}
            />
          </h1>
          <p className="hero__subtitle">{siteConfig.tagline}</p>
          <div className={styles.buttons}>
            <Link
              className={classnames(
                'button button--primary button--lg',
                styles.getStarted
              )}
              to={useBaseUrl('docs/introduction')}
            >
              Get Started
            </Link>
          </div>
        </div>
      </header>
      <main>
        {features && features.length && (
          <section className={classnames('padding-vert--xl', styles.features)}>
            <div className="container">
              <div className="row">
                {features.map((props, idx) => (
                  <Feature key={idx} {...props} />
                ))}
              </div>
            </div>
          </section>
        )}
      </main>
    </Layout>
  )
}

export default Home
