import query from './query.graphql'
import ProductSubnav from 'components/subnav'
import HashiStackMenu from '@hashicorp/react-hashi-stack-menu'
import Footer from 'components/footer'
import { open } from '@hashicorp/react-consent-manager'

export default function StandardLayout(props: Props): React.ReactElement {
  const { useCaseNavItems } = props.data

  return (
    <>
      <HashiStackMenu />
      <ProductSubnav
        menuItems={[
          { text: 'Overview', url: '/' },
          {
            text: 'Use Cases',
            submenu: useCaseNavItems,
          },
          {
            text: 'Enterprise',
            url: 'https://www.hashicorp.com/products/vault/enterprise',
          },
          'divider',
          { text: 'Tutorials', url: 'https://learn.hashicorp.com/vault' },
          { text: 'Docs', url: '/docs' },
          { text: 'API', url: '/api-docs' },
          { text: 'Community', url: '/community' },
        ]}
      />
      {props.children}
      <Footer openConsentManager={open} />
    </>
  )
}

StandardLayout.rivetParams = {
  query,
  dependencies: [],
}

interface Props {
  children: React.ReactChildren
  data: {
    useCaseNavItems: Array<{ slug: string; text: string }>
  }
}
