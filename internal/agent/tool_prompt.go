package agent

const toolPlannerPrompt = `Decide whether Luna should use a local tool before answering.

Available tools:
- web_search input: {"query": string, "max_results": int, "domains": []string}
- web_fetch input: {"url": string, "format": "markdown"}

Use tools for:
- Current or latest tax/compliance information.
- Official form names, filing instructions, portal manuals, notifications, circulars, and due dates.
- Questions about ITR, GST, TDS/TRACES, MCA/ROC/LLP filings, AIS/TIS/26AS, or official portal workflows.

Rules:
- Prefer official domains such as incometax.gov.in, gst.gov.in, mca.gov.in, cbic.gov.in, cbdt.gov.in, tdscpc.gov.in, and protean-tinpan.com.
- Use any conversation context provided below only when it matters for the current request.
- Use web_search first, then web_fetch when you need the actual page content.
- Do not invent URLs or unsupported tools.
- If the existing tool transcript is enough, return a final answer.

Respond with exactly one block and nothing else:

<tool_call>
{"tool":"web_search","input":{"query":"...","max_results":5}}
</tool_call>

or

<tool_call>
{"tool":"web_fetch","input":{"url":"https://...","format":"markdown"}}
</tool_call>

or

<final>
your answer
</final>`
