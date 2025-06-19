# Settings Configuration

## Overview

Notedown Planner stores its configuration settings in a YAML file on your local system. This document explains where settings are stored and what configuration options are available.

## Storage Location

The configuration file is stored at:

- **macOS/Linux**: `~/.notedown/planner/config.yaml`
- **Windows**: `%USERPROFILE%\.notedown\planner\config.yaml`

Where `~` or `%USERPROFILE%` represents your home directory.

## Configuration File Format

The configuration is stored in YAML format with the following structure:

```yaml
workspace_root: /path/to/your/workspace
```

## Configuration Options

### workspace_root
- **Type**: String (file path)
- **Required**: Yes
- **Description**: The root directory where Notedown Planner will store and manage your notes and planning documents
- **Example**: `/Users/username/Documents/Notedown` or `C:\Users\username\Documents\Notedown`

## Directory Structure

When you first run Notedown Planner, it will:

1. Create the configuration directory at `~/.notedown/planner/` if it doesn't exist
2. Create the `config.yaml` file after you complete the initial setup
3. Validate that your selected workspace directory is writable

## Manual Configuration

While it's recommended to use the application's settings interface, you can manually edit the configuration file if needed:

1. Close Notedown Planner
2. Navigate to the configuration file location
3. Edit `config.yaml` with any text editor
4. Save the file
5. Restart Notedown Planner

**Note**: Be careful when manually editing the configuration. Invalid paths or malformed YAML will prevent the application from starting properly.

## Troubleshooting

### Application asks for setup on every launch
- Check that the configuration file exists at the expected location
- Verify the `workspace_root` path in the config file is valid and accessible
- Ensure you have write permissions to both the config directory and workspace directory

### Configuration changes not saving
- Verify you have write permissions to `~/.notedown/planner/`
- Check that the disk isn't full
- Look for any error messages in the application logs

### Resetting configuration
To reset your configuration and start fresh:
1. Close Notedown Planner
2. Delete or rename the `config.yaml` file
3. Restart the application to go through setup again