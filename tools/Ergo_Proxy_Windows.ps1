#Change proxy address script
set-itemproperty -path "hkcu:Software\Microsoft\Windows\CurrentVersion\Internet Settings" -name AutoConfigURL -value "http://127.0.0.1:2000/proxy.pac" -type string

#For to apply changes for the Win Proxy
refresh-system


# The registry changes aren't seen until system is notified about it.
# Without this function you need to open Internet Settings window for changes to take effect. See http://goo.gl/OIQ4W4
function refresh-system() {
    $signature = @'
[DllImport("wininet.dll", SetLastError = true, CharSet=CharSet.Auto)]
public static extern bool InternetSetOption(IntPtr hInternet, int dwOption, IntPtr lpBuffer, int dwBufferLength);
'@

    $INTERNET_OPTION_SETTINGS_CHANGED   = 39
    $INTERNET_OPTION_REFRESH            = 37
    $type = Add-Type -MemberDefinition $signature -Name wininet -Namespace pinvoke -PassThru
    $a = $type::InternetSetOption(0, $INTERNET_OPTION_SETTINGS_CHANGED, 0, 0)
    $b = $type::InternetSetOption(0, $INTERNET_OPTION_REFRESH, 0, 0)
    return $a -and $b
}
