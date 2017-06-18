package layout

func defaultHeader() (string) {
  return `
<div class="grid">
<div class="grid-item 2/4"><span class="text-muted">Gravity</span><span style="class:#eee;">Build System</span></div>
  <div class="grid-item 2/4"></div>
</div>
<div class="grid">
  <div class="grid-item 1/3">
  </div>
  <div class="grid-item 2/3">
    <form class="form-inline">
      <input placeholder="Search project files" class="4/5">
      <button type="submit">Search</button>
    </form>
  </div>
</div>
`}

func DefaultLayout() (string) {
  return `
<html>` +
  // [HEADER] Define the default meta tags and pass in the CSS
 `<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Gravity</title>
    <style>` +
    CSS() + CSSIcons() + CSSExp() +
    `</style>
  </head>
  <body>
  <div class="container">` +
    // [Content]
    `<header>` +
      defaultHeader() +
    `</header>
      <main class="container">` +
          searchContent() +
      `</main>
    </div>
  </body>
</html>
`}
