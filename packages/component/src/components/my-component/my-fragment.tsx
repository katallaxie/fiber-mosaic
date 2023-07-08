import { Component, Prop, h } from '@stencil/core'

@Component({
  tag: 'my-fragment',
  shadow: true,
})
export class MyFragment {
  @Prop() deferred: boolean

  render() {
    return <div>A new fragment</div>
  }
}
