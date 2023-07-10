import { html, LitElement } from 'lit'
import { property, state } from 'lit/decorators.js'
import { until } from 'lit/directives/until.js'
import { unsafeHTML } from 'lit/directives/unsafe-html.js'

export class AppFragment extends LitElement {
  @property({ type: Boolean }) deferred = true

  @property({ type: String }) src = 'http://localhost:3000/fragment1'

  @property({ type: String }) method = 'GET'

  @state()
  private content = fetch(this.src)
    .then(r => r.text())
    .then(h => html`<div>${unsafeHTML(h)}</div>`)

  render() {
    return html`${until(this.content, html`<span>Loading...</span>`)}`
  }
}
