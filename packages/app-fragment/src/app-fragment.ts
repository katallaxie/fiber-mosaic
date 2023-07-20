import { html, LitElement } from 'lit'
import { property, customElement, state } from 'lit/decorators.js'
import { unsafeHTML } from 'lit/directives/unsafe-html.js'

@customElement('app-fragment')
export class AppFragment extends LitElement {
  /**
   * If true, the fragment will be loaded only when it is visible.
   * @attr deferred
   * @type {boolean}
   * @default false
   */
  @property()
  deferred: boolean = false

  /**
   * The URL of the fragment.
   * @attr src
   * @type {string}
   * @default ''
   */
  @property()
  src: string = ''

  /**
   * The content of the fragment.
   * @attr content
   * @type {string}
   * @default ''
   */
  @state()
  protected content = ''

  /**
   * The HTTP method to use when fetching the fragment.
   * @attr method
   * @type {string}
   * @default 'GET'
   * @see https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods
   * @see https://fetch.spec.whatwg.org/#methods
   * @see https://developer.mozilla.org/en-US/docs/Web/API/Request/method
   */
  @property()
  method = 'GET'

  /**
   * Loads the fragment.
   * @returns {Promise<void>}
   * @private
   * @see https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch
   * @see https://developer.mozilla.org/en-US/docs/Web/API/Request
   * @see https://developer.mozilla.org/en-US/docs/Web/API/Response
   * @see https://developer.mozilla.org/en-US/docs/Web/API/Body
   * @see https://developer.mozilla.org/en-US/docs/Web/API/Body/text
   */
  async load() {
    const response = await fetch(this.src, { method: this.method })
    const html = await response.text()
    this.content = html
  }

  /**
   * Lifecycle method that is called when the element is first updated.
   * @returns {Promise<void>}
   * @protected
   * @see https://lit.dev/docs/components/lifecycle/#update
   */
  async connectedCallback() {
    super.connectedCallback()
    await this.load()
  }

  /**
   * Renders the fragment.
   *
   * @returns {TemplateResult}
   */
  render() {
    if (!this.deferred) return html`<slot></slot>`

    return html`<div>${unsafeHTML(this.content)}</div>`
  }
}
