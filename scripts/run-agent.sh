#!/bin/bash
# Usage: ./scripts/run-{agent}.sh {feature_name}
# This script runs an agent in a loop until it reports completion.

AGENT_NAME="{{AGENT_NAME}}"
AGENT_NAME_LOWER="{{AGENT_NAME_LOWER}}"
FEATURE_NAME="$1"

if [ -z "$FEATURE_NAME" ]; then
  echo "Usage: $0 <feature_name>"
  echo "Example: $0 login"
  exit 1
fi

cd "$(dirname "$0")/.."

echo "========================================"
echo "Starting ${AGENT_NAME} agent"
echo "Feature: ${FEATURE_NAME}"
echo "========================================"
echo ""

ITERATION=1
MAX_ITERATIONS={{MAX_ITERATIONS}}

while true; do
  echo "[Iteration ${ITERATION}/${MAX_ITERATIONS}] Running ${AGENT_NAME} agent..."

  # Run the agent command
  echo "/${AGENT_NAME_LOWER} ${FEATURE_NAME}" | claude --print

  # Check for completion message
  LATEST_MSG=$(ls -t doc/agents/messages/${AGENT_NAME_LOWER}_to_orchestrator_*.md 2>/dev/null | head -1)

  if [ -n "$LATEST_MSG" ]; then
    if grep -q "Status: completed" "$LATEST_MSG"; then
      echo ""
      echo "========================================"
      echo "${AGENT_NAME} agent completed all tasks!"
      echo "Completion message: ${LATEST_MSG}"
      echo "========================================"
      break
    fi
  fi

  # Check max iterations
  if [ "$ITERATION" -ge "$MAX_ITERATIONS" ]; then
    echo ""
    echo "========================================"
    echo "Max iterations (${MAX_ITERATIONS}) reached."
    echo "${AGENT_NAME} agent stopped."
    echo "Run /workflow-adapter:validate to check progress."
    echo "========================================"
    break
  fi

  ITERATION=$((ITERATION + 1))

  echo ""
  echo "Waiting 10 seconds before next iteration..."
  echo "(Press Ctrl+C to stop)"
  echo ""
  sleep 10
done

echo ""
echo "${AGENT_NAME} agent finished."
