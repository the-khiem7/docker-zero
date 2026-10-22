# Saves a screenshot of one open Docker Zero window, as the built exe draws
# it, for the inventory's real references. Text whose accessible name matches
# -Hide (for example account paths) is painted over before saving.
#
#   powershell -File design/ui-inventory/real.ps1 -Window "Docker Zero settings" -Out design/ui-inventory/real/real-settings.png -Hide '\\Users\\'
#
# Window titles: "Docker Zero" (popup), "Docker Zero status" (corner box),
# "Docker Zero settings", "Set up Docker Zero".
param(
    [Parameter(Mandatory)][string]$Window,
    [Parameter(Mandatory)][string]$Out,
    [string]$Hide
)
Add-Type -AssemblyName UIAutomationClient, UIAutomationTypes, System.Drawing
$A = [System.Windows.Automation.AutomationElement]
$cond = New-Object System.Windows.Automation.PropertyCondition($A::NameProperty, $Window)
$win = $A::RootElement.FindFirst([System.Windows.Automation.TreeScope]::Children, $cond)
if (-not $win) { throw "window '$Window' is not open" }

if (-not ('InventoryShot' -as [type])) {
    Add-Type -TypeDefinition @'
using System; using System.Runtime.InteropServices;
public static class InventoryShot {
  [DllImport("user32.dll")] public static extern bool PrintWindow(IntPtr h, IntPtr dc, uint flags);
}
'@
}
$r = $win.Current.BoundingRectangle
$bmp = New-Object System.Drawing.Bitmap([int]$r.Width, [int]$r.Height)
$g = [System.Drawing.Graphics]::FromImage($bmp)
$dc = $g.GetHdc()
# PW_RENDERFULLCONTENT captures WebView2 content even when the window is covered.
[void][InventoryShot]::PrintWindow([IntPtr]$win.Current.NativeWindowHandle, $dc, 2)
$g.ReleaseHdc($dc)

if ($Hide) {
    $all = $win.FindAll([System.Windows.Automation.TreeScope]::Descendants, [System.Windows.Automation.Condition]::TrueCondition)
    $font = New-Object System.Drawing.Font('Segoe UI', 9)
    foreach ($e in $all) {
        if ($e.Current.Name -notmatch $Hide) { continue }
        $b = $e.Current.BoundingRectangle
        if ($b.IsEmpty) { continue }
        $x = [int]($b.X - $r.X); $y = [int]($b.Y - $r.Y)
        $bg = $bmp.GetPixel([Math]::Max(0, $x - 2), [Math]::Max(0, $y - 2))
        $g.FillRectangle((New-Object System.Drawing.SolidBrush($bg)), $x, $y, [int]$b.Width, [int]$b.Height)
        $g.DrawString('[hidden: account path]', $font, [System.Drawing.Brushes]::Gray, $x, $y)
    }
}
$g.Dispose()
New-Item -ItemType Directory -Force (Split-Path $Out) | Out-Null
$bmp.Save($Out, [System.Drawing.Imaging.ImageFormat]::Png)
"saved $([int]$r.Width)x$([int]$r.Height) to $Out"
