import { Component, Prop, h } from '@stencil/core'

@Component({
  tag: 'app-fragment',
  shadow: true,
})
export class AppFragment {
  /**
   * If `true`, the fragment will be fetched in the browser.
   */
  @Prop() deferred: boolean

  render() {
    return <div>A new fragment</div>
  }
}
