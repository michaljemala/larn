package main

/*
 * Here are the functions in this file:
 *
 * movemonst()		Routine to move the monsters toward the player
 * movemt(x,y)		Function to move a monster at (x,y) -- must determine where
 * mmove(x,y,xd,yd)	Function to actually perform the monster movement
 * movsphere() 		Function to look for and move spheres of annihilation
 */

/*
 * movemonst()		Routine to move the monsters toward the player
 *
 * This routine has the responsibility to determine which monsters are to
 * move, and call movemt() to do the move.
 * Returns no value.
 */
var w1, w1x, w1y [9]int
var tmp1, tmp2, tmp3, tmp4, distance int

func movemonst() {
	if c[TIMESTOP] != 0 {
		return /* no action if time is stopped */
	}
	if c[HASTESELF] != 0 {
		if c[HASTESELF]&1 == 0 {
			return
		}
	}
	if spheres != nil {
		movsphere() /* move the spheres of annihilation if any */
	}
	if c[HOLDMONST] != 0 {
		return /* no action if monsters are held */
	}

	if c[AGGRAVATE] != 0 { /* determine window of monsters to move */
		tmp1 = playery - 5
		tmp2 = playery + 6
		tmp3 = playerx - 10
		tmp4 = playerx + 11
		distance = 40 /* depth of intelligent monster movement */
	} else {
		tmp1 = playery - 3
		tmp2 = playery + 4
		tmp3 = playerx - 5
		tmp4 = playerx + 6
		distance = 17 /* depth of intelligent monster movement */
	}

	if level == 0 { /* if on outside level monsters can move in perimeter */
		if tmp1 < 0 {
			tmp1 = 0
		}
		if tmp2 > MAXY {
			tmp2 = MAXY
		}
		if tmp3 < 0 {
			tmp3 = 0
		}
		if tmp4 > MAXX {
			tmp4 = MAXX
		}
	} else { /* if in a dungeon monsters can't be on the perimeter (wall there) */
		if tmp1 < 1 {
			tmp1 = 1
		}
		if tmp2 > MAXY-1 {
			tmp2 = MAXY - 1
		}
		if tmp3 < 1 {
			tmp3 = 1
		}
		if tmp4 > MAXX-1 {
			tmp4 = MAXX - 1
		}
	}

	for j := tmp1; j < tmp2; j++ { /* now reset monster moved flags */
		for i := tmp3; i < tmp4; i++ {
			moved[i][j] = 0
		}
	}
	moved[lasthx][lasthy] = 0

	if c[AGGRAVATE] != 0 || c[STEALTH] == 0 { /* who gets moved? split for efficiency */
		for j := tmp1; j < tmp2; j++ { /* look thru all locations in window */
			for i := tmp3; i < tmp4; i++ {
				if mitem[i][j] != 0 { /* if there is a monster to move */
					if moved[i][j] == 0 { /* if it has not already been moved */
						movemt(i, j) /* go and move the monster */
					}
				}
			}
		}
	} else { /* not aggravated and not stealth */
		for j := tmp1; j < tmp2; j++ { /* look thru all locations in window */
			for i := tmp3; i < tmp4; i++ {
				if mitem[i][j] != 0 { /* if there is a monster to move */
					if moved[i][j] == 0 { /* if it has not already been moved */
						if stealth[i][j] != 0 { /* if it is asleep due to stealth */
							movemt(i, j) /* go and move the monster */
						}
					}
				}
			}
		}
	}

	if mitem[lasthx][lasthy] != 0 { /* now move monster last hit by player if not already moved */
		if moved[lasthx][lasthy] == 0 { /* if it has not already been moved */
			movemt(lasthx, lasthy)
			lasthx = w1x[0]
			lasthy = w1y[0]
		}
	}
}

/*
 * movemt(x,y)		Function to move a monster at (x,y) -- must determine where
 * 	int x,y;
 *
 * This routine is responsible for determining where one monster at (x,y) will
 * move to.  Enter with the monsters coordinates in (x,y).
 * Returns no value.
 */
var tmpitem int

func movemt(i, j int) {
	monst := mitem[i][j]
	switch monst { /* for half speed monsters */
	case TROGLODYTE, HOBGOBLIN, METAMORPH, XVART, INVISIBLESTALKER, ICELIZARD:
		if gltime&1 == 1 {
			return
		}
	}

	if c[SCAREMONST] != 0 { /* choose destination randomly if scared */
		xl := i + rnd(3) - 2
		if xl < 0 {
			xl = 0
		}
		if xl >= MAXX {
			xl = MAXX - 1
		}
		yl := j + rnd(3) - 2
		if yl < 0 {
			yl = 0
		}
		if yl >= MAXY {
			yl = MAXY - 1
		}
		tmp := item[xl][yl]
		if tmp != OWALL {
			if mitem[xl][yl] == 0 {
				if (mitem[i][j] != VAMPIRE) || (tmpitem != OMIRROR) {
					if tmp != OCLOSEDDOOR {
						mmove(i, j, xl, yl)
					}
				}
			}
		}
		return
	}
	if monster[monst].intelligence > 10-c[HARDGAME] { /* if smart monster */
		/* intelligent movement here -- first setup screen array */
		xl := tmp3 - 2
		yl := tmp1 - 2
		xh := tmp4 + 2
		yh := tmp2 + 2
		vxy(&xl, &yl)
		vxy(&xh, &yh)
		for k := yl; k < yh; k++ {
			for m := xl; m < xh; m++ {
				switch item[m][k] {
				case OWALL, OPIT, OTRAPARROW, ODARTRAP, OCLOSEDDOOR, OTRAPDOOR, OTELEPORTER:
					screen[m][k] = 127
					break
				case OMIRROR:
					if mitem[m][k] == VAMPIRE {
						screen[m][k] = 127
						break
					}
					fallthrough
				default:
					screen[m][k] = 0
				}
			}
		}
		screen[playerx][playery] = 1

		/*
		 * now perform proximity ripple from playerx,playery to
		 * monster
		 */
		xl = tmp3 - 1
		yl = tmp1 - 1
		xh = tmp4 + 1
		yh = tmp2 + 1
		vxy(&xl, &yl)
		vxy(&xh, &yh)
		var tmp byte
		for tmp = 1; tmp < byte(distance); tmp++ { /* only up to 20 squares away */
			for k := yl; k < yh; k++ {
				for m := xl; m < xh; m++ {
					if screen[m][k] == tmp { /* if find proximity n advance it */
						for z := 1; z < 9; z++ { /* go around in a circle */
							xtmp := m + diroffx[z]
							ytmp := k + diroffy[z]
							if screen[xtmp][ytmp] == 0 {
								screen[xtmp][ytmp] = tmp + 1
							}
							if xtmp == i && ytmp == j {
								goto out
							}
						}
					}
				}
			}
		}

	out:
		if tmp < byte(distance) { /* did find connectivity */
			/* now select lowest value around playerx,playery */
			for z := 1; z < 9; z++ { /* go around in a circle */
				xl := i + diroffx[z]
				yl := j + diroffy[z]
				if screen[xl][yl] == tmp {
					if mitem[xl][yl] == 0 {
						w1x[0] = xl
						w1y[0] = yl
						mmove(i, j, w1x[0], w1y[0])
						return
					}
				}
			}
		}
	}
	/* dumb monsters move here */
	xl := i - 1
	yl := j - 1
	xh := i + 2
	yh := j + 2
	if i < playerx {
		xl++
	} else if i > playerx {
		xh--
	}
	if j < playery {
		yl++
	} else if j > playery {
		yh--
	}
	for k := 0; k < 9; k++ {
		w1[k] = 10000
	}

	for k := xl; k < xh; k++ {
		for m := yl; m < yh; m++ { /* for each square compute distance to player */
			tmp := k - i + 4 + 3*(m-j)
			tmpitem = item[k][m]
			if tmpitem != OWALL || (k == playerx && m == playery) {
				if mitem[k][m] == 0 {
					if mitem[i][j] != VAMPIRE || tmpitem != OMIRROR {
						if tmpitem != OCLOSEDDOOR {
							w1[tmp] = (playerx-k)*(playerx-k) + (playery-m)*(playery-m)
							w1x[tmp] = k
							w1y[tmp] = m
						}
					}
				}
			}
		}
	}

	tmp := 0
	for k := 1; k < 9; k++ {
		if w1[tmp] > w1[k] {
			tmp = k
		}
	}

	if w1[tmp] < 10000 {
		if (i != w1x[tmp]) || (j != w1y[tmp]) {
			mmove(i, j, w1x[tmp], w1y[tmp])
		}
	}
}

/*
 * mmove(x,y,xd,yd)	Function to actually perform the monster movement
 * 	int x,y,xd,yd;
 *
 * Enter with the from coordinates in (x,y) and the destination coordinates
 * in (xd,yd).
 */
func mmove(aa, bb, cc, dd int) {
	var who string

	flag := 0 /* set to 1 if monster hit by arrow trap */
	if cc == playerx && dd == playery {
		hitplayer(aa, bb)
		moved[aa][bb] = 1
		return
	}
	i := item[cc][dd]
	if i == OPIT || i == OTRAPDOOR {
		switch mitem[aa][bb] {
		case SPIRITNAGA, PLATINUMDRAGON, WRAITH, VAMPIRE, SILVERDRAGON, POLTERGEIST, DEMONLORD, DEMONLORD + 1, DEMONLORD + 2, DEMONLORD + 3, DEMONLORD + 4, DEMONLORD + 5, DEMONLORD + 6, DEMONPRINCE:

		default:
			mitem[aa][bb] = 0 /* fell in a pit or trapdoor */
		}
	}
	tmp := mitem[aa][bb]
	mitem[cc][dd] = tmp
	if i == OANNIHILATION {
		if tmp >= DEMONLORD+3 { /* demons dispel spheres */
			cursors()
			lprintf("\nThe %s dispels the sphere!", monster[tmp].name)
			rmsphere(cc, dd) /* delete the sphere */
		} else {
			i, tmp, mitem[cc][dd] = 0, 0, 0
		}
	}
	stealth[cc][dd] = 1
	hitp[cc][dd] = hitp[aa][bb]
	if hitp[cc][dd] < 0 {
		hitp[cc][dd] = 1
	}
	mitem[aa][bb] = 0
	moved[cc][dd] = 1
	if tmp == LEPRECHAUN {
		switch i {
		case OGOLDPILE, OMAXGOLD, OKGOLD, ODGOLD, ODIAMOND, ORUBY, OEMERALD, OSAPPHIRE:
			item[cc][dd] = 0 /* leprechaun takes gold */
		}
	}

	if tmp == TROLL { /* if a troll regenerate him */
		if gltime&1 == 0 {
			if monster[tmp].hitpoints > hitp[cc][dd] {
				hitp[cc][dd]++
			}
		}
	}

	if i == OTRAPARROW { /* arrow hits monster */
		who = "An arrow"
		hitp[cc][dd] -= rnd(10) + level
		if hitp[cc][dd] <= 0 {
			mitem[cc][dd] = 0
			flag = 2
		} else {
			flag = 1
		}
	}
	if i == ODARTRAP { /* dart hits monster */
		who = "A dart"
		hitp[cc][dd] -= rnd(6)
		if hitp[cc][dd] <= 0 {
			mitem[cc][dd] = 0
			flag = 2
		} else {
			flag = 1
		}
	}
	if i == OTELEPORTER { /* monster hits teleport trap */
		flag = 3
		fillmonst(mitem[cc][dd])
		mitem[cc][dd] = 0
	}
	if c[BLINDCOUNT] != 0 {
		return /* if blind don't show where monsters are	 */
	}
	if know[cc][dd] {
		if flag != 0 {
			cursors()
		}
		switch flag {
		case 1:
			lprintf("\n%s hits the %s", who, monster[tmp].name)
			beep()
		case 2:
			lprintf("\n%s hits and kills the %s",
				who, monster[tmp].name)
			beep()
		case 3:
			lprintf("\nThe %s gets teleported", monster[tmp].name)
			beep()
		}
	}
	/*
	 * if yrepcount > 1 { know[aa][bb] &= 2;  know[cc][dd] &= 2; return
	 * }
	 */
	if know[aa][bb] {
		show1cell(aa, bb)
	}
	if know[cc][dd] {
		show1cell(cc, dd)
	}
}

/*
 * movsphere() 	Function to look for and move spheres of annihilation
 *
 * This function works on the sphere linked list, first duplicating the list
 * (the act of moving changes the list), then processing each sphere in order
 * to move it.  They eat anything in their way, including stairs, volcanic
 * shafts, potions, etc, except for upper level demons, who can dispel
 * spheres.
 * No value is returned.
 */
const SPHMAX = 20 /* maximum number of spheres movsphere can handle */
func movsphere() {
	var sph [SPHMAX]sphere

	// first duplicate sphere list
	x := 0
	for sp2 := spheres; sp2 != nil; sp2 = sp2.p { // look through sphere list
		if sp2.lev == level { // only if this level
			sph[x] = *sp2
			sph[x].p = nil // copy the struct
			x++
			if x > 1 {
				sph[x-2].p = &sph[x-1] // link pointers
			}
		}
	}
	if x == 0 {
		// no spheres
		return
	}

	for sp := &sph[0]; sp != nil; sp = sp.p { // look through sphere list
		x := sp.x
		y := sp.y
		if item[x][y] != OANNIHILATION {
			continue // not really there
		}
		sp.lifetime--
		if sp.lifetime < 0 { // has sphere run out of gas?
			rmsphere(x, y) // delete sphere
			continue
		}
		switch rnd(max(7, c[INTELLIGENCE]>>1)) { // time to move the sphere
		case 1:
		case 2: // change direction to a random one
			sp.dir = rnd(8)
		default: // move in normal direction
			dir := sp.dir
			len := sp.lifetime
			rmsphere(x, y)
			newsphere(x+diroffx[dir], y+diroffy[dir], dir, len)
		}
	}
}
