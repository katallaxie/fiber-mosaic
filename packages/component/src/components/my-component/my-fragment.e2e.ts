import { newSpecPage } from '@stencil/core/testing'
import { MyFragment } from './my-fragment'

describe('my-component', () => {
  it('renders', async () => {
    const { root } = await newSpecPage({
      components: [MyFragment],
      html: '<my-fragment></my-fragment>',
    })
    expect(root).toEqualHtml(`
      <my-fragment>
        <mock:shadow-root>
          <div>A new fragment</div>
        </mock:shadow-root>
      </my-fragment>
    `)
  })
})
