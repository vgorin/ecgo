// Copyright 2013-2014 Vasiliy Gorin.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

// Original Jerasure C/C++ code –
// Copyright 2007 James S. Plank
// See copyright notice inside *.c, *.h files

/*
 * Auxiliary routines required for Go integration
 *   added by vgorin
 */

// sizeofint returns sizeof(int); used by Go app to determine int <-> C.int compatability
int sizeofint();
