# Define color codes
$Red = [System.ConsoleColor]::Red
$Green = [System.ConsoleColor]::Green
$Yellow = [System.ConsoleColor]::Yellow
$NoColor = [System.ConsoleColor]::White

# Define the installation directory and binary name
$INSTALL_DIR = "$HOME\.local\kthcloud\bin"
$BINARY_NAME = "kthcloud"
$GITHUB_REPO = "Phillezi/kthcloud-cli"

# Detect OS and architecture
$OS = [System.Environment]::OSVersion.Platform
$ARCH = if ([System.Environment]::Is64BitOperatingSystem) { "x64" } else { "x86" }

# Determine OS and ARCH for the binary
# Powershell can be installed on other OS:s than Windows now, so make sure its running on windows 
switch ($OS) {
    "Win32NT" { $OS = "windows" }
    default { Write-Host "Unsupported OS: $OS" -ForegroundColor $Red; exit 1 }
}

switch ($ARCH) {
    "X64" { $ARCH = "amd64" }
    default { Write-Host "Unsupported architecture: $ARCH" -ForegroundColor $Red; exit 1 }
}

# Construct the download URL for the binary
$BINARY_URL = "https://github.com/$GITHUB_REPO/releases/latest/download/${BINARY_NAME}_${ARCH}_${OS}.exe"

# Function to show a loading spinner
function Show-Spinner {
    param (
        [string]$Url,
        [System.Net.WebClient]$WebClient
    )
    $spinstr = '|/-\'
    $delay = 100
    while ($WebClient.IsBusy) {
        foreach ($char in $spinstr.ToCharArray()) {
            Write-Host -NoNewline "$char" -ForegroundColor $Green
            Start-Sleep -Milliseconds $delay
            Write-Host -NoNewline "`b"
        }
    }
    Write-Host -NoNewline " " * $spinstr.Length -ForegroundColor $Green
}

# Create the install directory if it doesn't exist
if (-not (Test-Path $INSTALL_DIR)) {
    New-Item -Path $INSTALL_DIR -ItemType Directory | Out-Null
}

# Download the binary
Write-Host "Downloading $BINARY_NAME for $OS $ARCH..." -ForegroundColor $Green
$webClient = New-Object System.Net.WebClient
$downloadTask = $webClient.DownloadFileTaskAsync($BINARY_URL, "$INSTALL_DIR\$BINARY_NAME.exe")
Show-Spinner -Url $BINARY_URL -WebClient $webClient

# Wait for the download task to complete and handle errors
try {
    $downloadTask.Wait()
} catch {
    Write-Host "Failed to download the binary... :(" -ForegroundColor $Red
    Write-Host "Check if the URL is correct:"
    Write-Host $BINARY_URL
    exit 1
}

# Add binary path to the user PATH environment variable
$pathVariable = [System.Environment]::GetEnvironmentVariable("PATH", [System.EnvironmentVariableTarget]::User)
if (-not ($pathVariable -like "*$INSTALL_DIR*")) {
    [System.Environment]::SetEnvironmentVariable("PATH", "$pathVariable;$INSTALL_DIR", [System.EnvironmentVariableTarget]::User)
    Write-Host "Added $INSTALL_DIR to PATH environment variable" -ForegroundColor $Green
    Write-Host "Please restart PowerShell or open a new terminal to apply the changes."
} else {
    Write-Host "Path $INSTALL_DIR already added to PATH environment variable"
}

Write-Host "$BINARY_NAME installed successfully!" -ForegroundColor $Green
