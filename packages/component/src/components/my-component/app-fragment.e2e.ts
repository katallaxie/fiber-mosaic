import { newSpecPage } from '@stencil/core/testing'
import { AppFragment } from './app-fragment'

describe('my-component', () => {
  it('renders', async () => {
    const { root } = await newSpecPage({
      components: [AppFragment],
      html: '<app-fragment></app-fragment>',
    })
    expect(root).toEqualHtml(`
      <app-fragment>
        <mock:shadow-root>
          <div>A new fragment</div>
        </mock:shadow-root>
      </app-fragment>
    `)
  })
})
