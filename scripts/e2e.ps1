$ErrorActionPreference = 'Stop'

function Write-Step([string]$msg) { Write-Host ("[STEP] " + $msg) -ForegroundColor Cyan }

$base = "http://localhost:8080"

Write-Step "Health"
Invoke-RestMethod -Method GET "$base/health" | Out-Host

Write-Step "Admin login"
$loginBody = @{ email='admin@university.edu'; password='admin123' } | ConvertTo-Json
$adminResp = Invoke-RestMethod -Method POST "$base/api/v1/auth/login" -Headers @{ 'Content-Type'='application/json' } -Body $loginBody
$adminToken = $adminResp.token
if (-not $adminToken) { throw 'Admin login failed' }

Write-Step "Register student"
$unique = [guid]::NewGuid().ToString('N').Substring(0,8)
$regBody = @{ name='John Doe'; email="john+$unique@student.edu"; password='password123'; role='student'; dept='Computer Science' } | ConvertTo-Json
$regResp = Invoke-RestMethod -Method POST "$base/api/v1/auth/register" -Headers @{ 'Content-Type'='application/json' } -Body $regBody
$studentToken = $regResp.token
$studentId = $regResp.user.id
if (-not $studentToken -or -not $studentId) { throw 'Student registration failed' }

Write-Step "Student applies for leave"
$applyBody = @{ leave_type='Medical'; reason='Doctor appointment'; start_date=(Get-Date).ToString('s')+'Z'; end_date=(Get-Date).AddDays(1).ToString('s')+'Z' } | ConvertTo-Json
$applyResp = Invoke-RestMethod -Method POST "$base/api/v1/leaves/apply" -Headers @{ 'Content-Type'='application/json'; 'Authorization'="Bearer $studentToken" } -Body $applyBody
$leaveId = $applyResp.data.id
if (-not $leaveId) { throw 'Leave apply failed' }

Write-Step "Student: my leaves"
Invoke-RestMethod -Method GET "$base/api/v1/leaves/my" -Headers @{ 'Authorization'="Bearer $studentToken" } | Out-Host

Write-Step "Admin: pending leaves"
$pending = Invoke-RestMethod -Method GET "$base/api/v1/leaves/pending" -Headers @{ 'Authorization'="Bearer $adminToken" }
($pending | ConvertTo-Json -Depth 5) | Out-Host
$approveId = if ($pending.data -and $pending.data.Count -gt 0) { $pending.data[0].id } else { $leaveId }

Write-Step "Admin: approve leave $approveId"
$approveBody = @{ approved = $true; remarks = 'Approved' } | ConvertTo-Json
Invoke-RestMethod -Method PUT "$base/api/v1/leaves/$approveId/approve" -Headers @{ 'Content-Type'='application/json'; 'Authorization'="Bearer $adminToken" } -Body $approveBody | Out-Host

Write-Step "Admin: mark attendance for student $studentId"
$markBody = @{ student_id = $studentId; date = (Get-Date).ToString('s')+'Z'; present = $true } | ConvertTo-Json
Invoke-RestMethod -Method POST "$base/api/v1/attendance/mark" -Headers @{ 'Content-Type'='application/json'; 'Authorization'="Bearer $adminToken" } -Body $markBody | Out-Host

Write-Step "Student: my attendance"
Invoke-RestMethod -Method GET "$base/api/v1/attendance/my" -Headers @{ 'Authorization'="Bearer $studentToken" } | Out-Host

Write-Step "Admin: attendance by student id"
$start = (Get-Date).AddDays(-7).ToString('s')+'Z'
$end = (Get-Date).ToString('s')+'Z'
Invoke-RestMethod -Method GET "$base/api/v1/attendance/student/$studentId?start_date=$start&end_date=$end" -Headers @{ 'Authorization'="Bearer $adminToken" } | Out-Host

Write-Step "Student: notifications"
Invoke-RestMethod -Method GET "$base/api/v1/notifications/my" -Headers @{ 'Authorization'="Bearer $studentToken" } | Out-Host
Invoke-RestMethod -Method GET "$base/api/v1/notifications/unread-count" -Headers @{ 'Authorization'="Bearer $studentToken" } | Out-Host

Write-Step "Admin: users list"
Invoke-RestMethod -Method GET "$base/api/v1/users?page=1&limit=10" -Headers @{ 'Authorization'="Bearer $adminToken" } | Out-Host

Write-Step "Admin: analytics dashboard"
Invoke-RestMethod -Method GET "$base/api/v1/analytics/dashboard" -Headers @{ 'Authorization'="Bearer $adminToken" } | Out-Host

Write-Host "E2E test completed successfully." -ForegroundColor Green



