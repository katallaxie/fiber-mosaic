import { html, LitElement } from 'lit'
import { property, state } from 'lit/decorators.js'
import { until } from 'lit/directives/until.js'
import { unsafeHTML } from 'lit/directives/unsafe-html.js'

export class AppFragment extends LitElement {
  /**
   * If true, the fragment will be loaded only when it is visible.
   * @attr deferred
   * @type {boolean}
   * @default false
   */
  @property({ type: Boolean }) deferred = false

  /**
   * The URL of the fragment.
   * @attr src
   * @type {string}
   * @default ''
   */
  @property({ type: String }) src = ''

  /**
   * The HTTP method to use when fetching the fragment.
   * @attr method
   * @type {string}
   * @default 'GET'
   * @see https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods
   * @see https://fetch.spec.whatwg.org/#methods
   * @see https://developer.mozilla.org/en-US/docs/Web/API/Request/method
   */
  @property({ type: String }) method = 'GET'

  @state()
  private content = fetch(this.src)
    .then(r => r.text())
    .then(h => html`<div>${unsafeHTML(h)}</div>`)

  /**
   * Renders the fragment.
   *
   * @returns {TemplateResult}
   */
  render() {
    if (!this.deferred) return html`<slot></slot>`

    return html`${until(this.content, html`<slot></slot>`)}`
  }
}
