/*
 * *****************************************************************************
 * Copyright (C) 2019-2024 Chrystian Huot <chrystian@huot.qc.ca>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>
 * ****************************************************************************
 */

import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { MatSelectChange } from '@angular/material/select';
import { RdioScannerAdminService, Group, Tag } from '../../../admin.service';

@Component({
    selector: 'rdio-scanner-admin-talkgroup',
    templateUrl: './talkgroup.component.html',
})
export class RdioScannerAdminTalkgroupComponent {
    @Input() form: FormGroup | undefined;

    @Output() blacklist = new EventEmitter<void>();

    @Output() remove = new EventEmitter<void>();

    leds: string[];

    get alerts(): string[] {
        return Object.keys(this.adminService.Alerts || {});
    }

    get groups(): Group[] {
        return this.form?.root.get('groups')?.value as Group[];
    }

    get tags(): Tag[] {
        return this.form?.root.get('tags')?.value as Tag[];
    }

    constructor(private adminService: RdioScannerAdminService) {
        this.leds = this.adminService.getLeds();
    }

    async playAlert(event: MatSelectChange): Promise<void> {
        if (event.value) await this.adminService.playAlert(event.value);
    }
}
