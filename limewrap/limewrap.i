/* limewrap.i */
%module limewrap
%{
#include <lime/LimeSuite.h>
%}

%insert(cgo_comment_typedefs) %{
#cgo LDFLAGS: -lLimeSuite
%}

#define _DOXYGEN_ONLY_

%include "/usr/include/lime/LimeSuite.h"