# Slow command notifier

Reduce the risk of getting distracted by showing an OS notification after a slow terminal command completes.

## Using

### Use with fish shell

```fish
function fish_prompt
    # MUST be the first thing in fish_prompt!
    set -l cmd_status $status

    show_slow_command_notif \
        -cmd=$history[1] \
        -duration=$CMD_DURATION \
        -status=$cmd_status \
        -prev_fg_app_asn=$_FISH_PROMPT_FG_APP_ASN \
        -curr_fg_app_asn=(lsappinfo front)

    # ... rest of your prompt
end

set -g _FISH_PROMPT_FG_APP_ASN (lsappinfo front)
function preexec --on-event fish_preexec
  set -g _FISH_PROMPT_FG_APP_ASN (lsappinfo front)
end
```
