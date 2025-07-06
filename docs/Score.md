## Scoring table:

<head>
  <style>
    table {
      width: 100%;
      border-collapse: collapse;
      border-radius: 8px;
      overflow: hidden;
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    }
    thead {
      <!-- background-color: #4b5563; -->
      color: white;
    }
    th {
      padding: 12px 15px;
      text-align: left;
      font-weight: 600;
    }
    td {
      padding: 12px 15px;
      border-bottom: 1px solid #e5e7eb;
    }
    tr:last-child td {
      border-bottom: none;
    }
    tr:nth-child(even) {
      <!-- background-color: #f9fafb; -->
    }
    tr:hover {
      background-color: #f3f4f6;
    }
    .badge {
      display: inline-block;
      padding: 4px 8px;
      border-radius: 9999px;
      font-weight: 500;
      text-align: center;
      white-space: nowrap;
      color: white;
    }
    .badge-blue { background-color: #3B82F6; }
    .badge-gray { background-color: #6B7280; }
    .badge-green { background-color: #10B981; }
    .badge-purple { background-color: #8B5CF6; }
    .badge-orange { background-color: #F97316; }
    .badge-red { background-color: #EF4444; }
    .badge-yellow { background-color: #F59E0B; }
    .score {
      display: inline-block;
      width: 60px;
      padding: 4px 4px;
      border-radius: 9999px;
      font-weight: 500;
      text-align: center;
      color: white;
    }
    .score-good { background-color: #10B981; }
    .score-medium { background-color: #F59E0B; }
    .score-poor { background-color: #EF4444; }
    .model-name {
      font-weight: 500;
      color: #1f2937;
    }
    .eval-data {
      font-family: monospace;
      color: #4b5563;
    }
    .detail-text {
      color: #6b7280;
      font-size: 0.9em;
      max-width: 300px;
    }
  </style>
</head>
<body>
  <table>
    <thead>
      <tr>
        <th>Name</th>
        <th>Algorithm</th>
        <th>Eval</th>
        <th>Score</th>
        <th>Detail</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td class="model-name">name</td>
        <td><span class="badge badge-blue">Algo</span></td>
        <td class="eval-data">MAE: , Loss: </td>
        <td><span>0.xxx</span></td>
        <td class="detail-text"></td>
      </tr>
    </tbody>
  </table>
</body>
