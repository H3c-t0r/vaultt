import Hero from '@hashicorp/react-hero'
import SectionHeader from '@hashicorp/react-section-header'
import UseCases from '@hashicorp/react-use-cases'
import TextSplitWithCode from '@hashicorp/react-text-split-with-code'
import TextSplitWithLogoGrid from '@hashicorp/react-text-split-with-logo-grid'
import Button from '@hashicorp/react-button'

import BeforeAfterDiagram from '../../components/before-after-diagram'

export default function Homepage() {
  return (
    <div id="page-home">
      <div className="g-section-block page-wrap">
        {/* Hero */}
        <Hero
          data={{
            backgroundImage: {
              alt: null,
              size: 55617,
              url:
                'https://www.datocms-assets.com/2885/1539894412-vault-bg.jpg',
              width: 3196,
            },
            backgroundTheme: 'light',
            buttons: [
              {
                external: false,
                title: 'Download',
                url: 'https://www.vaultproject.io/downloads',
              },
              {
                external: false,
                title: 'Get Started with Vault',
                url: 'https://www.vaultproject.io/intro/getting-started',
              },
            ],
            centered: false,
            description:
              'Secure, store and tightly control access to tokens, passwords, certificates, encryption keys for protecting secrets and other sensitive data using a UI, CLI, or HTTP API.',
            theme: 'vault-gray',
            title: 'Manage Secrets and Protect Sensitive Data',
            videos: [
              {
                name: 'UI',
                playbackRate: 2,
                src: [
                  {
                    srcType: 'mp4',
                    url:
                      'https://www.datocms-assets.com/2885/1543956852-vault-v1-0-ui-opt.mp4',
                  },
                ],
              },
              {
                name: 'CLI',
                playbackRate: 2,
                src: [
                  {
                    srcType: 'mp4',
                    url:
                      'https://www.datocms-assets.com/2885/1543956847-vault-v1-0-cli-opt.mp4',
                  },
                ],
              },
            ],
          }}
        />

        {/* Text Section */}

        <section className="g-grid-container remove-bottom-padding">
          <SectionHeader
            headline="Secure dynamic infrastructure across clouds and environments"
            description="The shift from static, on-premise infrastructure to dynamic, multi-provider infrastructure changes the approach to security. Security in static infrastructure relies on dedicated servers, static IP addresses, and a clear network perimeter. Security in dynamic infrastructure is defined by ephemeral applications and servers, trusted sources of user and application identity, and software-based encryption."
          />
        </section>

        {/* Before-After Diagram */}

        <section className="g-grid-container before-after">
          <BeforeAfterDiagram
            beforeImage={{
              url:
                'https://www.datocms-assets.com/2885/1579635889-static-infrastructure.svg',
              format: 'svg',
            }}
            beforeHeadline="Static Infrastructure"
            beforeContent={`Datacenters with inherently high-trust networks with clear network perimeters.

#### Traditional Approach

- High trust networks
- A clear network perimeter
- Security enforced by IP Address`}
            afterImage={{
              url:
                'https://www.datocms-assets.com/2885/1579635892-dynamic-infrastructure.svg',
              format: 'svg',
            }}
            afterHeadline="Dynamic Infrastructure"
            afterContent={`Multiple clouds and private datacenters without a clear network perimeter.

#### Vault Approach


- Low-trust networks in public clouds
- Unknown network perimeter across clouds
- Security enforced by Identity`}
          />
        </section>

        {/* Use cases */}

        <section>
          <div className="g-grid-container">
            <UseCases
              theme="vault"
              items={[
                {
                  title: 'Secrets Management',
                  description:
                    'Audit access, automatically Centrally store, access, and deploy secrets across applications, systems, and infrastructure',
                  image: {
                    alt: null,
                    format: 'png',
                    url:
                      'https://www.datocms-assets.com/2885/1575422126-secrets.png',
                  },
                  link: {
                    external: false,
                    title: 'Learn more',
                    url: '/use-cases/secrets-management',
                  },
                },
                {
                  title: 'Data Encryption',
                  description:
                    'Keep secrets and application data secure with one centralized workflow to encrypt data in flight and at rest',
                  image: {
                    alt: null,
                    format: 'png',
                    url:
                      'https://www.datocms-assets.com/2885/1575422166-encryption.png',
                  },
                  link: {
                    external: false,
                    title: 'Learn more',
                    url: '/use-cases/data-encryption',
                  },
                },
                {
                  title: 'Identity-based Access',
                  description:
                    'Authenticate and access different clouds, systems, and endpoints using trusted identities',
                  image: {
                    alt: null,
                    format: 'png',
                    url:
                      'https://www.datocms-assets.com/2885/1575422201-identity.png',
                  },
                  link: {
                    external: false,
                    title: 'Learn more',
                    url: '/use-cases/identity-based-access',
                  },
                },
              ]}
            />
          </div>
        </section>
        {/* Principles / Text & Content Blocks */}

        <section className="no-spacing">
          <div className="g-grid-container">
            <SectionHeader headline="Vault Principles" />
          </div>

          <TextSplitWithCode
            textSplit={{
              heading: 'API-driven',
              content:
                'Use policy to codify, protect, and automate access to secrets',
            }}
            codeBlock={{
              options: { showWindowBar: true },
              code:
                '$ curl \n\t--header "X-Vault-Token: ..." \n\t--request POST \n\t--data @payload.json \n\thttps://127.0.0.1:8200/v1/secret/config',
            }}
          />

          <TextSplitWithLogoGrid
            textSplit={{
              heading: 'Identity Plugins',
              content: 'Seamlessly integrate any trusted identity provider',
              textSide: 'right',
            }}
            logoGrid={[
              'aws',
              'microsoft-azure',
              'google',
              {
                url:
                  'https://www.datocms-assets.com/2885/1556657783-oktalogo.svg',
              },
              { url: 'https://www.datocms-assets.com/2885/1539817287-cf.svg' },
              'alibaba-cloud',
              {
                url: 'https://www.datocms-assets.com/2885/1506540149-black.svg',
              },
              'kubernetes',
              'github',
            ]}
          />

          <TextSplitWithLogoGrid
            textSplit={{
              heading: 'Extend and integrate',
              content:
                'Securely manage secrets and access through a centralized workflow',
            }}
            logoGrid={[
              'mysql',
              'cassandra',
              'oracle',
              'aws',
              'mongodb',
              {
                url:
                  'https://www.datocms-assets.com/2885/1508434209-consul_primarylogo_fullcolor.svg',
              },
              {
                url:
                  'https://www.datocms-assets.com/2885/1539817686-microsoft-sql-server.svg',
              },
              'postgresql',
              'microsoft-azure',
            ]}
          />
        </section>

        <section className="g-grid-container">
          <SectionHeader
            headline="Open Source and Enterprise"
            description="Vault Open Source addresses the technical complexity of managing secrets by leveraging trusted identities across distributed infrastructure and clouds. Vault Enterprise addresses the organizational complexity of large user bases and compliance requirements with collaboration and governance features."
          />
          <div className="button-container">
            <Button
              title="Learn More"
              url="https://www.hashicorp.com/products/vault/enterprise"
            />
          </div>
        </section>
      </div>
    </div>
  )
}
