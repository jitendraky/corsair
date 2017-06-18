package layout

import ()

// TODO
// 1. Santize query before printing in template
func searchContent() (string) {
  return `
<h5>
  <span class="text-muted">{{ .FoundCount }}</span> results searching for '<span class="text-muted">{{ .Query }}</span>'
</h5>
<div class="grid">
  <div class="grid-item 3/5"></div>
  <div class="grid-item 1/5"><span class="text-muted" style="font-size:12px;">sort</span></div>
  <div class="grid-item 1/5"><span style="font-size:10px;margin-top:4px;">a-to-z</span></div>
</div>
<table>
  {{ range $path, $result := .FoundFiles }}
    <thead>
      <tr>
        <th>Path</th>
        <th>{{ $path }}</th>
      </tr>
    </thead>
    {{ range $s, $e := $result }}
     <tr>
        <th>Term</th>
      </tr>
      <tr>
        <td><strong>Start</strong></td>
        <td>Line</td>
        <td>{{ $e.Start.Line }}</td>
        <td>Col</td>
        <td>{{ $e.Start.Col }}</td>
      </tr>
      <tr>
        <td><strong>End</strong></td>
        <td>Line</td>
        <td>{{ $e.End.Line }}</td>
        <td>Col</td>
        <td>{{ $e.End.Col }}</td>
      </tr>
      <tr>
        <th>Snippet</th>
     </tr>
      <tr>
        <td><strong>Start</strong></td>
        <td>Line</td>
       <td>{{ $e.Snippet.Start }}</td>
       <td>End</td>
       <td>{{ $e.Snippet.End }}</td>
      </tr>
      <tr>
        <td>{{ $e.Snippet.Text }}</td>
      </tr>
    {{ end  }}
  {{ end  }}
</table>
  }`
}
