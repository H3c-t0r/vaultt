import * as React from 'react'
import { Products } from '@hashicorp/platform-product-meta'
import classNames from 'classnames'
import Image from 'next/image'
import Button from '@hashicorp/react-button'
import s from './style.module.css'

interface IoUsecaseSectionsProps {
  brand?: Products | 'neutral'
  sections: [IoUsecaseSectionProps, IoUsecaseSectionProps]
}

export default function IoUseCaseSestions({
  brand = 'neutral',
  sections,
}: IoUsecaseSectionsProps): React.ReactElement {
  return (
    <>
      {sections.map((section, index) => {
        // Index is stable
        // eslint-disable-next-line react/no-array-index-key
        return <IoUsecaseSection key={index} brand={brand} {...section} />
      })}
    </>
  )
}

interface IoUsecaseSectionProps {
  brand?: Products | 'neutral'
  eyebrow: string
  heading: string
  description: React.ReactNode
  media?: {
    src: string
    width: string
    height: string
    alt: string
  }
  cta?: {
    text: string
    link: string
  }
}

function IoUsecaseSection({
  brand = 'neutral',
  eyebrow,
  heading,
  description,
  media,
  cta,
}: IoUsecaseSectionProps): React.ReactElement {
  return (
    <section className={classNames(s.section, s[brand])}>
      <div className={s.container}>
        <p className={s.eyebrow}>{eyebrow}</p>
        <div className={s.columns}>
          <div className={s.column}>
            <h2 className={s.heading}>{heading}</h2>
            {media ? <p className={s.description}>{description}</p> : null}
            {cta ? (
              <div className={s.cta}>
                <Button
                  title={cta.text}
                  url={cta.link}
                  theme={{
                    brand: brand,
                  }}
                />
              </div>
            ) : null}
          </div>
          <div className={s.column}>
            {media ? (
              // eslint-disable-next-line jsx-a11y/alt-text
              <Image {...media} />
            ) : (
              <div className={s.description}>{description}</div>
            )}
          </div>
        </div>
      </div>
    </section>
  )
}
